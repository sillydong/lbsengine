import axios from "axios";
import {MessageBox,Indicator} from "mint-ui";
import router from "../router/index"
import * as qs from "qs";

// 创建axios实例
const service = axios.create({
  baseURL: process.env.BASE_API, // api的base_url
  timeout: 5000,                  // 请求超时时间
});

// request拦截器
service.interceptors.request.use((config) => {
  Indicator.open();
  // Do something before request is sent
  config.headers.post['Content-Type'] = 'application/x-www-form-urlencoded';
  if (config.method === 'post') {
    config.data = qs.stringify(config.data);
  }
  return config;
}, (error) => {
  // Do something with request error
  console.log(error); // for debug
  MessageBox.alert("请求失败5",error);
  return Promise.reject(error);
});

// respone拦截器
service.interceptors.response.use(
    (response) => {
      Indicator.close()
      console.log(response.data)
      const status = response.data.status;
      if (status == 1) {
        return response
      } else {
        const err = response.data.error;
        switch (typeof(err)) {
          case "string":
            MessageBox.alert("请求失败4", err);
            break;
          case "object":
            MessageBox.alert("请求失败3", Object.values(err).join());
            break;
          case "array":
            MessageBox.alert("请求失败2", Object.values(err).join());
            break;
        }
        return Promise.reject(response);
      }
    }, error => {
      Indicator.close();
      console.log('err' + error);// for debug
      MessageBox.alert("请求失败1", error.message);
      return Promise.reject(error);
    }
);

export function get(url, params) {
  return new Promise((resolve, reject) => {
    service.get(url, {params: params}).then(response => {
      resolve(response.data);
    }, err => {
      reject(err)
    }).catch((error) => {
      reject(error)
    });
  })
}

export function post(url, params) {
  return new Promise((resolve, reject) => {
    service.post(url, params).then(response => {
      resolve(response.data);
    }, err => {
      reject(err)
    }).catch((error) => {
      reject(error)
    });
  })
}

export function del(url){
  return new Promise((resolve,reject)=>{
    service.del(url).then(response => {
      resolve(response.data);
    },err=>{
      reject(err)
    }).catch(err => {
      reject(err)
    });
  })
}
