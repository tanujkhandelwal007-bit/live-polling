# Live Polling Frontend

This project is the frontend for a real-time polling application built with React and Vite.

## Project Structure

```bash
frontend/
├── index.html
├── package.json
├── vite.config.js
├── public/
├── src/
│   ├── App.css
│   ├── App.jsx
│   ├── index.css
│   ├── main.jsx
│   └── assets/
├── ws-test.html
├── README.md
├── .gitignore
├── .oxlintrc.json
├── node_modules/
└── package-lock.json
```

## Main Features

- User registration and login
- JWT-based authentication
- Create polls with multiple options
- View poll details from a URL like `/poll/:id`
- Vote on a poll
- Live result updates
- Close polls as the owner
- Copy shareable link
- Real-time communication through WebSocket

## Tech Stack

- React
- Vite
- JavaScript
- CSS
- WebSocket

## Setup

Install dependencies:

```bash
npm install
```

Run the app in development mode:

```bash
npm run dev
```

The frontend expects the backend to run at:

```bash
http://localhost:8080
```

## Expected Backend APIs

This frontend communicates with a backend using endpoints such as:

- `POST /auth/register`
- `POST /auth/login`
- `POST /polls`
- `GET /polls/:id`
- `GET /polls/:id/results`
- `POST /polls/:id/vote`
- `PUT /polls/:id/close`
- WebSocket: `ws://localhost:8080/polls/:id/ws`

## App Flow

1. User registers or logs in.
2. User creates a poll with a question and multiple options.
3. The app sends poll data to the backend and redirects to the poll page.
4. Poll participants vote and see real-time count updates.
5. The creator can close the poll when needed.

## Notes

- JWT is stored in `localStorage`.
- Poll state is managed in the frontend based on the current URL.
- Real-time result updates happen through a WebSocket connection.

## Future Improvements

- Add proper error handling and validation for backend errors
- Improve UI/UX for polling results
- Add user-specific poll history
- Add production-ready deployment configuration
