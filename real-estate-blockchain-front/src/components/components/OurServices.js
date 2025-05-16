// src/components/OurServices.js
import React, { useEffect } from 'react';
import './OurServices.css';
import AOS from 'aos';
import 'aos/dist/aos.css';

const services = [
  {
    title: '블록체인 매물 인증',
    desc: '모든 거래 기록을 블록체인에 저장하여 위변조를 방지합니다.',
    icon: '🔗',
  },
  {
    title: '실시간 위치 기반 검색',
    desc: '지도 기반으로 원하는 지역의 매물을 한눈에 확인할 수 있습니다.',
    icon: '📍',
  },
  {
    title: '스마트 계약 기반 거래',
    desc: '계약 프로세스를 자동화하여 거래 과정을 투명하고 간단하게 만듭니다.',
    icon: '📄',
  },
];

const OurServices = () => {
  useEffect(() => {
    AOS.init({ duration: 1000 });
  }, []);

  return (
    <section className="our-services" id="services">
      <h2 data-aos="fade-up">💼 Our Services</h2>
      <p className="subtext" data-aos="fade-up" data-aos-delay="100">
        블록체인 기술을 기반으로 신뢰할 수 있는 부동산 거래 환경을 제공합니다.
      </p>
      <div className="service-cards">
        {services.map((s, idx) => (
          <div className="service-card" key={idx} data-aos="fade-up" data-aos-delay={200 + idx * 100}>
            <div className="icon">{s.icon}</div>
            <h3>{s.title}</h3>
            <p>{s.desc}</p>
          </div>
        ))}
      </div>
    </section>
  );
};

export default OurServices;
