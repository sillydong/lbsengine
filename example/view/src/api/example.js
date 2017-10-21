import {get,post} from "@/utils/fetch";

export function api(params){
  return get('/api/query',params)
}
