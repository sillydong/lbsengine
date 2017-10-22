package main

import (
	"database/sql"
	"github.com/Code-Hex/echo-static"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/go-gorp/gorp"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sillydong/lbsengine/engine"
	"github.com/sillydong/lbsengine/types"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	//tidb
	dbconfig := &mysql.Config{
		User:   "root",
		Passwd: "",
		DBName: "lbsexample",
		Net:    "tcp",
		Addr:   "127.0.0.1:4000",
	}
	db, err := sql.Open("mysql", dbconfig.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	client := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}

	client.AddTableWithName(PoiData{}, "poi_data").SetKeys(true, "id")

	//engine
	eg := &engine.Engine{}
	eg.Init(nil)

	//echo
	addr := ":8877"
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: ResourceSkipper,
	}))
	e.Use(middleware.Recover())

	e.Use(static.ServeRoot("/", &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "view/dist",
	}))
	api := e.Group("/api")
	api.POST("/add", func(context echo.Context) error {
		//解析对象
		item := new(PoiData)
		if err := context.Bind(item); err != nil {
			return context.JSON(http.StatusInternalServerError, Response{
				Status: 0,
				Error:  "请求的参数不正确",
			})
		}

		//写tidb
		err := client.Insert(item)
		if err != nil {
			context.Logger().Error(err)
			return context.JSON(http.StatusInternalServerError, Response{
				Status: 0,
				Error:  err.Error(),
			})
		}
		if item.Id > 0 {
			//写engine
			eg.Add(&types.IndexedDocument{
				DocId:     uint64(item.Id),
				Latitude:  item.Latitude,
				Longitude: item.Longitude,
				Fields:    nil,
			})
			return context.JSON(http.StatusOK, Response{
				Status: 1,
				Data:   "添加成功",
			})
		} else {
			return context.JSON(http.StatusInternalServerError, Response{
				Status: 0,
				Error:  "数据库写入失败",
			})
		}

	})
	api.DELETE("/del/:id", func(context echo.Context) error {
		id := context.Param("id")
		fid, _ := strconv.Atoi(id)
		//删tidb
		_, err := client.Delete(&PoiData{Id: int64(fid)})
		if err != nil {
			context.Logger().Error(err)
		}

		//删engine
		eg.Remove(uint64(fid))

		return context.JSON(http.StatusOK, Response{
			Status: 1,
			Data:   "删除成功",
		})
	})
	api.GET("/query", func(context echo.Context) error {
		latitude := context.QueryParam("latitude")
		longitude := context.QueryParam("longitude")
		offset := context.QueryParam("offset")

		flat, _ := strconv.ParseFloat(latitude, 10)
		flng, _ := strconv.ParseFloat(longitude, 10)
		foff, _ := strconv.Atoi(offset)

		//查engine
		resp := eg.Search(&types.SearchRequest{
			Latitude:  flat,
			Longitude: flng,
			CountOnly: false,
			Offset:    foff,
			Limit:     10,
		})
		if resp.Count > 0 && resp.Docs != nil {
			//拼数据
			ids := make([]interface{}, len(resp.Docs))
			for i, doc := range resp.Docs {
				ids[i] = doc.DocId
				//context.Logger().Printf("%+v\n",doc)
			}
			var datas []PoiData
			_, err := client.Select(&datas, "select * from poi_data where id in (?"+strings.Repeat(",?", len(resp.Docs)-1)+")", ids...)
			if err != nil {
				context.Logger().Error(err)
			}

			datasmap := make(map[uint64]PoiData, len(datas))
			for _, data := range datas {
				datasmap[uint64(data.Id)] = data
			}
			for _, doc := range resp.Docs {
				doc.Model = datasmap[doc.DocId]
			}

			return context.JSON(http.StatusOK, Response{
				Status: 1,
				Data:   resp,
			})
		} else {
			return context.JSON(http.StatusOK, Response{
				Status: 0,
				Error:  "无结果",
			})
		}
	})
	api.GET("/init", func(context echo.Context) error {
		var datas []PoiData
		_, err := client.Select(&datas, "select * from poi_data")
		if err != nil {
			context.Logger().Error(err)
		}

		context.Logger().Printf("about to import %v lines", len(datas))

		for _, item := range datas {
			context.Logger().Print(item.Id)
			eg.Add(&types.IndexedDocument{
				DocId:     uint64(item.Id),
				Latitude:  item.Latitude,
				Longitude: item.Longitude,
				Fields:    nil,
			})
		}

		return context.JSON(http.StatusOK, "导入完成")
	})
	e.Logger.Fatal(e.Start(addr))
}

func ResourceSkipper(c echo.Context) bool {
	if strings.HasPrefix(c.Path(), "/static") {
		return true
	}
	return false
}

type PoiData struct {
	Id         int64   `db:"id" form:"-"`
	Name       string  `db:"name" form:"name"`
	AmapId     string  `db:"amapid" form:"amapid"`
	Location   string  `db:"location" form:"location"`
	Latitude   float64 `db:"latitude" form:"latitude"`
	Longitude  float64 `db:"longitude" form:"longitude"`
	CreateTime string  `db:"create_time" form:"create_time"`
}

type Response struct {
	Status int         `json:"status"`
	Error  string      `json:"error"`
	Data   interface{} `json:"data"`
}
