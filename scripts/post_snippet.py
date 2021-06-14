#!/usr/bin/env python3

import requests
import sys

def get_token():
    try:
        with open(".token") as f:
            return f.read().strip()
    except:
        return None

def build_headers():
    token = get_token()
    if token:
        return { "Authorization": token }
    return {}

def main():
    if len(sys.argv) != 2:
        print(f"Usage: {sys.argv[0]} <snippet>")
        sys.exit(1)
    language = input("language: ")
    with open(sys.argv[1]) as f:
        contents = f.read()
    r = requests.post(f"http://localhost:5000/api/v1/snippets", headers=build_headers(), json={
        "Language": language,
        "Contents": contents
    })
    print(r.text)

if __name__ == "__main__":
    main()
