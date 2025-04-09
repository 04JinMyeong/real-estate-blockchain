// routes/authRoutes.js
const express = require('express');
const bcrypt = require('bcrypt');
const jwt = require('jsonwebtoken');
const router = express.Router();

const users = []; // 임시 저장소 (DB 대신)

router.post('/signup', async (req, res) => {
  const { email, password } = req.body;
  const exists = users.find(user => user.email === email);
  if (exists) return res.status(409).json({ message: '이미 존재하는 이메일입니다' });

  const hashedPassword = await bcrypt.hash(password, 10);
  users.push({ email, password: hashedPassword });
  res.json({ message: '회원가입 완료' });
});

router.post('/login', async (req, res) => {
  const { email, password } = req.body;
  const user = users.find(user => user.email === email);
  if (!user) return res.status(400).json({ message: '사용자 없음' });

  const valid = await bcrypt.compare(password, user.password);
  if (!valid) return res.status(401).json({ message: '비밀번호 틀림' });

  const token = jwt.sign({ email }, process.env.JWT_SECRET, { expiresIn: '1h' });
  res.json({ message: '로그인 성공', token, email: user.email  }); // ✅ 이렇게 명확히 user.email 사용!
});

module.exports = router;
