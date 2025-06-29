import express from 'express';
import dotenv from 'dotenv';
import pool, { testDBConnection } from './db/dbconn.js';

dotenv.config();

const app = express();
const port = process.env.PORT || 3000;

// Test DB connection, then start the server
const startServer = async () => {
  try {
    await testDBConnection();
    app.listen(port, () => {
      console.log(`Server started at http://localhost:${port}`);
    });
  } catch (err) {
    console.error('Failed to connect to DB:', err);
  }
};

startServer();


// Define a basic route
app.get('/', (req, res) => {
  res.send('Hello from Express + PostgreSQL!');
});

startServer();
