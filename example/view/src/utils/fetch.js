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
  MessageBox.alert("请求失败",error);
  return Promise.reject(error);
});

// respone拦截器
service.interceptors.response.use(
    (response) => {
      Indicator.close()
      console.log(response.data)
      const status = response.data.status;
      // 50014:Token 过期了 50012:其他客户端登录了 50008:非法的token
      if (status === "0") {
        return response
      } else {
        const err = response.data.error;
        switch (typeof(err)) {
          case "string":
            MessageBox.alert("请求失败", err);
            break;
          case "object":
            MessageBox.alert("请求失败", Object.values(err).join());
            break;
          case "array":
            MessageBox.alert("请求失败", Object.values(err).join());
            break;
        }
        return Promise.reject(response);
      }
    }, error => {
      Indicator.close();
      console.log('err' + error);// for debug
      MessageBox.alert("请求失败", error.message);
      if(error.response.status===401){
        store.dispatch('Logout').then(() => {
          router.push({path: '/login'})
        });
      }
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
