
# Library Management System

I have developed a comprehensive Library Management System using React.js for the frontend and Golang for the backend. The system incorporates robust authentication and authorization mechanisms, supporting three distinct user roles: User, Admin, and Owner.


## Features

Admin Features
- Onboard Readers: Admins can onboard new readers and manage their profiles.
- Book Management: Admins can:
- Add new books to the library collection.
- Delete books from the library.
- Update book details, including title, author, and availability.
- Approve or Disapprove book requests made by readers.
- CRUD Operations: Admins have full Create, Read, Update, and Delete capabilities for book records.

Reader Features
- Book Search: Readers can search for books within the library.
- QR Code Generation: Readers can generate a QR code for a selected book, which includes detailed information about the book.
- Book Requests: Readers can submit requests to issue or return books.



## Technical Implementation

**Authentication & Authorization:** Implemented secure user authentication with password hashing (64-bit encryption) and JSON Web Tokens (JWT) for session management. Tokens are stored in the browser’s local storage and expire after 24 hours.

**Form Validation:** All forms are equipped with thorough validation mechanisms to ensure data integrity.

**Responsive Design:** The frontend features fully responsive React components, ensuring a seamless user experience across various devices.

**Notifications:** Integrated the React Toastify library for effective and user-friendly notifications.


## Installation

Install my-project with npm

Run locally on localhost:3000 (frontend)

On the same time, clone backend repo and run it on localhost:8000 (backend server)

Frontend setup :

```bash
   git clone https://github.com/GujralNaman/library-management-frontend.git

```

```bash
  npm install my-project
  cd my-project
```

```bash
npm start
```

Backend setup :

```bash
git clone https://github.com/GujralNaman/library-management-backend.git
```

```bash
cd your_project_name
```

```bash
go run main.go
```
    