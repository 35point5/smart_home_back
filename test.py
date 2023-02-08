import time

import requests

t = time.time()
cnt = 0
for i in range(100):
    time.sleep(0.1)
    resp = requests.post("http://mogician.cc/smart_home/api/user/ping")
    if resp.status_code == 200:
        cnt = cnt + 1
print(cnt, "request success in", time.time() - t, "seconds.")
