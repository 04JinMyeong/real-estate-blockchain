// server.js
require('dotenv').config();

const express = require('express');
const cors = require('cors');
const dotenv = require('dotenv');
const authRoutes = require('../auth-backend/routes/authRoutes');

dotenv.config(); // .env 로드

const app = express();
app.use(cors());
app.use(express.json());

app.use('/api/auth', authRoutes);

const PORT = process.env.PORT || 3001;
app.listen(PORT, () => {
  console.log(`✅ 서버 실행 중: http://localhost:${PORT}`);
});
