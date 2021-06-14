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
    snippet = input("snippet: ")
    r = requests.get(f"http://localhost:5000/api/v1/snippets/{snippet}/comments", headers=build_headers())
    print(r.text)

if __name__ == "__main__":
    main()
