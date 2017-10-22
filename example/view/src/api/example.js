import {get,post,del} from "@/utils/fetch";

export function api_add(form){
  return post('/api/add',form)
}

export function api_del(id){
  return del('/api/del/'+id)
}

export function api_query(params){
  return get("/api/query",params)
}

