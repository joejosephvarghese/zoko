# 🚀 Zoko Server

## 📜 Description

**Zoko** is a real-time messaging application built using **Go** and **Socket.IO**. It enables users to communicate instantly and effortlessly.

The backend handles communication and data processing using **Socket.IO**, ensuring seamless and responsive real-time messaging.

---

## 🌟 Features

📨 **Real-time messaging:** Users can send and receive messages instantly.

🔐 **User authentication:** Secure authentication ensures that only authorized users can access the application.

✍️ **Typing indicators:** Displays when another user is typing a message.

🟢 **Online presence:** Users can see the online status of other users.

📜 **Message history:** Stores and displays message history for easy reference.

---

## 🛠 Technologies Used

🐹 **Go:** A powerful programming language known for its efficiency and concurrency.

🔄 **Socket.IO:** Enables real-time, bidirectional communication between the server and clients.

---

## 📋 Prerequisites

Ensure that you have the following software installed on your system:

🟡 **Go** (v1.16 or above)

---

## ⚙️ Configuration

The application can be configured using the **.env.example** file in the **server** folder. Make sure to set up the necessary environment variables before starting the application.

---

## 🏁 Getting Started

To run the **Zoko Server** application on your local machine, follow these steps:

### 📥 Clone the Repository

```sh
git clone https://github.com/joejosephvarghese/zoko.git
```

### 🖥️ Server Side Setup

1. Navigate to the server directory:
   ```sh
   cd zoko/server
   ```
2. Install dependencies:
   ```sh
   make deps || go mod tidy
   ```
3. Run the Go backend:
   ```sh
   make run || go run ./cmd/api/main.go
   ```

---

## 📖 API Documentation

📄 **Swagger UI:**

Access the API documentation at:

🔗 [Swagger UI](http://localhost:8080/swagger/index.html)

---

## 🎯 Usage

🚀 Start the backend service and connect clients for real-time messaging.

---

## 🤝 Contributing

💡 Contributions to **Zoko Server** are welcome! If you find any bugs or want to suggest improvements, please **open an issue** or **submit a pull request** on the GitHub repository.

---

## 📜 License

📝 This project is licensed under the **MIT License**. Feel free to use, modify, and distribute the code as per the license terms.

---

## 🙌 Acknowledgments

🔹 **Go:** [https://golang.org/](https://golang.org/)

🔹 **Socket.IO:** [https://socket.io/](https://socket.io/)

---

## 📬 Contact

📧 **Joe Joseph Varghese** - [joejosephvarghese94@gmail.com](mailto:joejosephvarghese94@gmail.com)

🚀 We hope you enjoy using **Zoko Server**! Happy messaging! 🎉
