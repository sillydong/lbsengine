<template>
  <div>
    <div>
      <span>{{ label }}</span>
      <el-amap vid="amap" :plugin="plugin" class="amap-demo" :center="center">
      </el-amap>
    </div>
    <ul v-infinite-scroll="loadmore" infinite-scroll-disabled="loading"
        infinite-scroll-distance="10">
      <li v-for="item in list">
        <div class="line">
          <div>{{item.DocId}}.{{ item.Model.Name }}</div>
          <div>Lat:{{ item.Model.Latitude}} Lng:{{ item.Model.Longitude }}</div>
          <div>Distance: {{ item.Distance}}</div>
        </div>

      </li>
    </ul>
  </div>
</template>

<script>
  //import { mapState, mapActions, mapMutations } from 'vuex';
  import {api_query} from "@/api/example"

  export default {
    components: {},
    data() {
      let self = this;
      return {
        label: '正在定位...',
        loading: false,
        list: [],
        query: {
          latitude: 0.0,
          longitude: 0.0,
          offset: 0
        },
        center: [121.59996, 31.197646],
        plugin: [{
          pName: 'Geolocation',
          events: {
            init(o) {
              o.getCurrentPosition((status, result) => {
                if (result && result.position) {
                  console.log(result.position)
                  self.query.latitude = result.position.lat;
                  self.query.longitude = result.position.lng;

                  self.label = "location: lat:" + self.query.latitude + " lng:"
                      + self.query.longitude;

                  self.loadmore();
                  self.$nextTick();
                }
              })
            }
          }

        }]
      }
    },
    computed: {},
    methods: {
      loadmore() {
        if(this.query.latitude!==0.0 && this.query.longitude!==0.0){
          api_query(this.query).then(response => {
            let items = response.data.Docs
            if (items.length > 0) {
              for (let i = 0; i < items.length; i++) {
                this.list.push(items[i])
              }
              this.query.offset += 10
              this.loading = false
            } else {
              console.log("empty")
            }
          }).catch(err => {
            this.loading = false
          })
        }
      }
    },
    watch: {},
    directives: {},
    filters: {},
    created() {
    },
    mounted() {
    },
  }
</script>

<style scoped>
  .amap-demo {
    height:100px;
  }
  .line{
    text-align: left;
  }
  ul{
    -webkit-padding-start: 20px;
  }
</style>
