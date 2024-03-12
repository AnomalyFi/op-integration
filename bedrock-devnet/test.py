import time
import jwt
import requests
import logging
import devnet

log = logging.getLogger()

l1_rpc_url = 'http://10.153.238.11:8545'
jwt_secret = '0xfad2709d0bb03bf0e8ba3c99bea194575d3e98863133d1af638ed056d1d59345'

res = devnet.debug_dumpBlock_local(l1_rpc_url, jwt_secret)
print(res)
