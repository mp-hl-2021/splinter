#!/usr/bin/env python3

import requests
import sys

def main():
    username = input("username: ")
    password = input("password: ")
    r = requests.post("http://localhost:5000/api/v1/authenticate", json={
        "Username": username,
        "Password": password
    })
    print(r.text)
    try:
        token = r.json()["Token"]
        with open(".token", "w") as f:
            print(token, file=f)
    except Exception as e:
        print(e)
        sys.exit(1)

if __name__ == "__main__":
    main()
