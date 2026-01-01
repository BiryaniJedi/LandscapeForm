## Forms table
Basic for dev, will update with all fields later
```go
id         uuid() pk
created_by uuid() NN
created_at timestamptz NN
first_name text NN
last_name  text NN
home_phone character(10) NN
```
## Shrubs table
```go
form_id     uuid()
    References forms(id) ODC
num_shrubs  int NNCheck (num_shrubs >= 0)
```
## Pesticides table
```go
form_id         uuid()
    References forms(id) ODC
pesticide_name  string NN
```
# Routes
```
POST   /forms/pesticide  
POST   /forms/shrub  
GET    /forms  
GET    /forms/{id}  
PUT    /forms/{id}  
DELETE /forms/{id}  
GET    /forms/{id}/pdf  
POST   /forms/import/pdf  
```