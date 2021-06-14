#!/usr/bin/env python3

import requests

def main():
    username = input("username: ")
    password = input("password: ")
    r = requests.post("http://localhost:5000/api/v1/create_account", json={
        "Username": username,
        "Password": password
    })
    print(r.text)

if __name__ == "__main__":
    main()
