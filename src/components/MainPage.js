// src/pages/MainPage.js
import React, { useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import './MainPage.css';
import { FaBitcoin, FaUserShield, FaCalendarCheck } from 'react-icons/fa';
// import OurServices from './OurServices';
import AOS from 'aos';
import 'aos/dist/aos.css';
import Wave from './Wave'

const MainPage = ({ user }) => {
  const serviceRef = useRef(null);
  const digitalDesignRef = useRef(null);
  const awesomeSupportRef = useRef(null);
  const easyCustomizeRef = useRef(null);
  const navigate = useNavigate();

  useEffect(() => {
    AOS.init({ duration: 1000, once: false, mirror: true,offset: 60,
    easing: "ease-in-out", });
  }, []);

  const scrollToServices = () => {
    serviceRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const scrollToDigitalDesign = () => {
    digitalDesignRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const scrollToAwesomeSupport = () => {
    awesomeSupportRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const scrollToEasyCustomize = () => {
    easyCustomizeRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const handleStart = () => {
    navigate('/signup');
  };

  return (
    <div className="main-page-wrapper">
      <section className="hero-section">
        <Wave />  {/* 상단에 위치 */}

        <div className="hero-content-wrapper">
          <div className="hero-content-box">
            <h1>블록체인 기반의 부동산 플랫폼</h1>
            <p>신뢰할 수 있는 매물 등록과 거래 이력, AI 분석까지 제공합니다.</p>
            <div className="cta-buttons">
              <button className="btn-primary" onClick={handleStart}>시작하기</button>
              <button className="btn-secondary" onClick={scrollToServices}>더 알아보기</button>
            </div>
          </div>
        </div>
      </section>





      {/* Services Section */}
      <section className="services-section" ref={serviceRef}>
        <h2 className="services-title" data-aos="fade-up">Our Services</h2>
        <p className="services-sub" data-aos="fade-up" data-aos-delay="100">
          신뢰할 수 있는 기술과 경험으로 최고의 부동산 경험을 제공합니다.
        </p>
        <div className="services-grid">
          <div
            className="service-box"
            onClick={scrollToDigitalDesign}
            style={{ cursor: 'pointer' }}
            data-aos="fade-up"
            data-aos-delay="200"
          >
            <div className="icon-circle">
              <FaBitcoin size={28} />
            </div>
            <h3>Block-chain</h3>
            <p>블록체인기술로 안전한 거래</p>
          </div>

          <div
            className="service-box active"
            onClick={scrollToAwesomeSupport}
            style={{ cursor: 'pointer' }}
            data-aos="fade-up"
            data-aos-delay="300"
          >
            <div className="icon-circle">
              <FaUserShield size={28} />
            </div>
            <h3>DID</h3>
            <p>탈중앙화된 신원인증 기술</p>
          </div>

          <div
            className="service-box"
            onClick={scrollToEasyCustomize}
            style={{ cursor: 'pointer' }}
            data-aos="fade-up"
            data-aos-delay="400"
          >
            <div className="icon-circle">
              <FaCalendarCheck size={28} />
            </div>
            <h3>Reservation</h3>
            <p>예약기능 최적화</p>
          </div>
        </div>
      </section>

      {/* Detail Sections Wrapped in White Background */}
      <section className="service-detail-wrapper">
        <section className="service-detail" ref={digitalDesignRef}>
          <div className="detail-grid" data-aos="fade-right">
            <img src="/Blockchain1.png" alt="UI Design" />
            <div className="detail-text">
              <h2>Block-chain</h2>
              <p>TrueHome은 블록체인 기술을 활용하여 거래과정을 안전하게 보호합니다.
                위변조가 불가능하고 인증된 사용자만 자산에 접근할 수 있습니다.
              </p>
            </div>
          </div>
        </section>

        <section className="service-detail" ref={awesomeSupportRef}>
          <div className="detail-grid" data-aos="fade-left">
            <img src="/Did2.png" alt="Different illustration" />
            <div className="detail-text">
              <h2>DID</h2>
              <p>탈중앙화된 신원 인증을 통해 사용자와 자산을 연결합니다.
                관리자 승인 후 중재자가 개입하지 않는 등록이 이루어집니다.
              </p>
            </div>
          </div>
        </section>

        <section className="service-detail" ref={easyCustomizeRef}>
          <div className="detail-grid" data-aos="fade-right">
            <img src="/Reservation3.png" alt="Customize" />
            <div className="detail-text">
              <h2>Reservation</h2>
              <p>예약 기능을 통해 사용자는 원하는 매물을 손쉽게 확인하고 중개인과 일정을 예약할 수 있어 편리합니다.</p>
            </div>
          </div>
        </section>

      </section>
    </div>
  );
};

export default MainPage;
