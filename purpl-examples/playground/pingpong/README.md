# 🏓 pingpong: two clients, one server

| who? | what? |
| ----------- | ----------- |
| goodclient | sends request to server |
| badclient | sends request to server |
| server_two | sends a name, phone number, street, age and sex as response |
| interceptor | minimizes the response depending on the client JWT |

## 🧪 Try it
```
go run server/server_two.go
go run client/clients.go
```

## 🥸 Data minimization opions 
- reduction
- noising
- generalization

Original message from server:
```
Name: "Ken Guru", PhoneNumber: "+0123456789", Street: "Straße des 17 Juni", Age: 35, Sex: "male"
```

The reduced result look like this:
```
-------------------------
Message from server for goodclient: name:"Ken Guru"  phoneNumber:"+"  street:"Str"  age:20  sex:"male"
Message from server for badclient:  name:"K"  phoneNumber:"+"  street:"S"  age:31  sex:"m"
-------------------------
```


## 🔑 Use of JSON Web Tokens

Check and generate them here: [jwt.io](https://jwt.io/).

Our token's secret: ```none```.

Right now our JWTs look like this:

### 🧳 payload for goodToken:
```
{
 	"policy": {
 	  "allowed": {
 		"name": "string",
 		"sex": "string"
 	  },
 	  "generalized": {
 		"phoneNumber": "string"
 	  },
 	  "noised": {
 		"age": "int"
 	  },
 	  "reduced": {
 		"street": "string"
 	  }
 	},
 	"exp": 1688843806,
 	"iss": "test"
}
```

### 🧳 payload for badToken:
```
{
 	"policy": {
 	"allowed": {},
 	"generalized": {
 		"age": "int",
 		"name": "string",
 		"phoneNumber": "string",
 		"sex": "string",
 		"street": "string"
 	},
 	"noised": {},
   	"reduced": {}
 	},
 	"exp": 1688483421,
 	"iss": "test"
}
```


- The clients append their respective JWTs to their request's context.
- The server's gRPC interceptor compares the outgoing response fields with the JWT's ```allowed```, ```generalized```, ```noised``` and ```reduced``` data fields. Allowed fields will be left untouched. Minimzed fields will be minimzed. Unmentioned fields will be suppressed to 1 or nil

## 🧭 Rodamap
- ...
