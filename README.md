<h1 align="center">
  Synapsis Project
  <br>
</h1>

<h4 align="center"> Synapsis Project is a Online Store Application, users can make transaction in this app </h4>
       
## Build App & Database
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![MySQL](https://img.shields.io/badge/mysql-%2300f.svg?style=for-the-badge&logo=mysql&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-black?style=for-the-badge&logo=JSON%20web%20tokens)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![GitHub](https://img.shields.io/badge/github-%23121011.svg?style=for-the-badge&logo=github&logoColor=white)

## ⚙️ Installing and Running from Github

Installing and running the app from github repository <br>
To clone and run this application, you'll need [Git](https://git-scm.com) and [Golang](https://go.dev/dl/) installed on your computer. From your command line:

```bash
# Clone this repository
$ git clone https://github.com/vaniliacahya/synapsis-project.git

# Go into the repository
$ cd synapsis-project

# Run the app
$ go run main.go
```

## Feature

- Customers <br>
This feature consist endpoint register for new customer, and login as a customer
  
| Customer | Endpoint | Param | JWT Token | Function |
| --- | --- | --- | --- | --- |
| `POST` | /register | - | NO | Feature register for new customer |
| `POST` | /login | - | NO | Feature login customer |


- Products <br>
This feature consist endpoint view product
  
| Product | Endpoint | Param | JWT Token | Function |
| --- | --- | --- | --- | --- |
| `GET` | /products | - | NO | Feature for customer to view product |


- Carts <br>
This feature consist endpoint to add, update, delete or get cart customer
  
| Cart | Endpoint | Param | JWT Token | Function |
| --- | --- | --- | --- | --- |
| `POST` | /carts | - | YES | Feature for customer to add or update product in cart |
| `GET` | /carts | - | YES | Feature for customer to get carts |
| `DELETE` | /carts | /:id | YES | Feature for customer to delete product in cart |

- Order <br>
This feature consist endpoint order
  
| Order | Endpoint | Param | JWT Token | Function |
| --- | --- | --- | --- | --- |
| `POST` | /checkout | - | YES | Feature for customer to make an order |

## ERD
<img src="ERD.png">

## Authors

Vanilia Cahya Nugroho
       
  Reach me:

  ![Gmail](https://img.shields.io/badge/vanilia.cahya23@gmail.com-EA4335.svg?style=for-the-badge&logo=gmail&logoColor=white)
  [![GitHub](https://img.shields.io/badge/vaniliacahya-%23121011.svg?style=for-the-badge&logo=github&logoColor=white)](https://github.com/vaniliacahya)
  [![Linkedin](https://img.shields.io/badge/vaniliacahya-0A66C2.svg?style=for-the-badge&logo=linkedin&logoColor=white)](https://www.linkedin.com/in/vaniliacahya/)
  [![WhatsApp](https://img.shields.io/badge/vaniliacahya-25D366.svg?style=for-the-badge&logo=whatsapp&logoColor=white)](https://api.whatsapp.com/send/?phone=%2B6281249690397&text=Hello&type=phone_number&app_absent=0)

       
 <p align="right">(<a href="#top">back to top</a>)</p>
<h3>
<p align="center">:copyright: 2022 </p>
</h3>
<!-- end -->
