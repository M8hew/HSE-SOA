import time
import requests
import json

# WAIT FOR THE SERVER TO START
time.sleep(20)

# REGISTER USER
register_url = 'http://localhost:8080/register'
register_headers = {
    "Content-Type": "application/json"
}
register_data = {
    "username": "matthew",
    "password": "example_password"
}

response = requests.post(register_url, headers=register_headers, data=json.dumps(register_data))

if response.status_code == 200:
    response_json = response.json()
    token = response_json.get("token")
    if token:
        print(f"Token received: {token}")
    else:
        print("Token not found in the response.")
        exit(1)
else:
    print(f"Registration request failed with status code {response.status_code}")
    exit(1)

# WAIT 
time.sleep(3)

# CREATING POST
create_post_url = 'http://localhost:8080/posts'
create_post_headers = {
    "Authorization": f"Bearer {token}",
    "Content-Type": "application/json"
}
create_post_data = {
    "content": "This is a new post."
}


for i in range(5):
    response = requests.post(create_post_url, headers=create_post_headers, data=json.dumps(create_post_data))

    if response.status_code == 200:
        print("Post created successfully.")
        post_id = response_json.get("Id")
        break
    else:
        print(f"Create post request failed with status code {response.status_code}")
        time.sleep(3)
else:
    print("Failed to create post after 5 attempts.")
    exit(1)

# ADD REACTION
add_reaction_url = f'http://localhost:8080/posts/{post_id}/view'
post_headers = {
    "Authorization": f"Bearer {token}",
    "Content-Type": "application/json"
}

requests.post(add_reaction_url, headers=post_headers)


# GET STAT
get_stat_url = f'http://localhost:8080/posts/{post_id}/stats'
get_stat_headers = {
    "accept": "application/json"
}

response = requests.get(get_stat_url, headers=get_stat_headers)
if response.status_code != 200 or response.json().get("views") != 1:
    print("Failed to retrieve post stats")
    exit(1)
