local jwt = require "resty.jwt"
local validators = require "resty.jwt-validators"
local cjson = require "cjson"

if ngx.var.request_method ~= "OPTIONS" then
 local jwtToken = ngx.var.http_Authorization

 if jwtToken == nil then
  ngx.status = ngx.HTTP_UNAUTHORIZED
  ngx.header.content_type = "application/json; charset=utf-8"
  ngx.say("{\"error\": \"Forbidden\"}")
  ngx.exit(ngx.HTTP_UNAUTHORIZED)
 end
 local claim_spec = {
  exp = validators.is_not_expired()
 }

 --local _, _, strippedToken = string.find(jwtToken, "Bearer%s+(.+)")
 local jwt_obj = jwt:verify("D474A1E00FC53F8D19E2795CC452D2E490C147773C25EB96E1559A22F32AA201", jwtToken, claim_spec)
 if not jwt_obj["verified"] then
  ngx.status = ngx.HTTP_UNAUTHORIZED
  ngx.header.content_type = "application/json; charset=utf-8"
  ngx.say("{\"error\": \"INVALID_JWT\"}")
  ngx.exit(ngx.HTTP_UNAUTHORIZED)
 end

 local jwt_obj_raw = jwt:load_jwt(jwtToken)
 local encoded = cjson.encode(jwt_obj_raw)
 local decoded = cjson.decode(encoded)
 local emailOverride = decoded.payload.email

 ngx.req.read_body()
 local req = ngx.req.get_body_data()
 if req ~= nil then
  local decodedReq = cjson.decode(req)
  decodedReq.email = emailOverride
  encodedReq = cjson.encode(decodedReq)
  ngx.req.set_body_data(encodedReq)
 else
  ngx.var.requesterEmail = emailOverride
 end

end