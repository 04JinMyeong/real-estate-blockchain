// src/pages/MainPage.js
import React, { useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import './MainPage.css';
import { FiGrid, FiLayers, FiTool } from 'react-icons/fi';
import OurServices from './OurServices';
import AOS from 'aos';
import 'aos/dist/aos.css';
import Wave from './Wave'

const MainPage = () => {
  const serviceRef = useRef(null);
  const digitalDesignRef = useRef(null);
  const awesomeSupportRef = useRef(null);
  const easyCustomizeRef = useRef(null);
  const navigate = useNavigate();

  useEffect(() => {
    AOS.init({ duration: 1000, once: false });
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
  <div className="hero-content-wrapper">
    <div className="hero-content-box">
      <h1>블록체인 기반의 부동산 플랫폼</h1>
      <p>신뢰할 수 있는 매물 등록과 거래 이력, AI 분석까지 제공합니다.</p>
      <div className="cta-buttons">
        <button>시작하기</button>
        <button>더 알아보기</button>
      </div>
    </div>
  </div>
  <Wave />  {/* Wave를 Hero 안에서 하단에 위치 */}
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
              <FiGrid size={28} />
            </div>
            <h3>Digital Design</h3>
            <p>정교하고 아름다운 매물 UI 제공</p>
          </div>

          <div
            className="service-box active"
            onClick={scrollToAwesomeSupport}
            style={{ cursor: 'pointer' }}
            data-aos="fade-up"
            data-aos-delay="300"
          >
            <div className="icon-circle">
              <FiLayers size={28} />
            </div>
            <h3>Different</h3>
            <p>기존 부동산과 다른 블록체인 기반 매물 정보 공개</p>
          </div>

          <div
            className="service-box"
            onClick={scrollToEasyCustomize}
            style={{ cursor: 'pointer' }}
            data-aos="fade-up"
            data-aos-delay="400"
          >
            <div className="icon-circle">
              <FiTool size={28} />
            </div>
            <h3>Easy to customize</h3>
            <p>관리자 및 중개인을 위한 맞춤형 대시보드</p>
          </div>
        </div>
      </section>

      {/* Detail Sections Wrapped in White Background */}
      <section className="service-detail-wrapper">
        <section className="service-detail" ref={digitalDesignRef}>
          <div className="detail-grid">
            <img src="/images/ui-design.jpg" alt="UI Design" />
            <div className="detail-text">
              <h2>Digital Design</h2>
              <p>
                TrueHome은 사용자가 쉽게 매물을 탐색하고 등록할 수 있도록 직관적인 사용자 인터페이스(UI)를 제공합니다.
                다양한 필터 기능과 반응형 디자인으로 모든 디바이스에서 편리하게 사용 가능합니다.
              </p>
            </div>
          </div>
        </section>

        <section className="service-detail" ref={awesomeSupportRef}>
          <div className="detail-grid">
            <img src="/images/transparent.jpg" alt="Transparency" />
            <div className="detail-text">
              <h2>Different</h2>
              <p>
                기존 부동산 서비스는 정보의 비대칭성과 위변조 가능성이 존재합니다. TrueHome은 블록체인 기술을 활용해
                모든 매물 정보, 등록 이력, 거래 기록을 누구나 열람 가능한 형태로 제공합니다.
                이러한 투명성은 신뢰를 기반으로 한 부동산 거래 문화를 만들어갑니다.
              </p>
            </div>
          </div>
        </section>

        <section className="service-detail" ref={easyCustomizeRef}>
          <div className="detail-grid">
            <img src="/images/customize.jpg" alt="Customize" />
            <div className="detail-text">
              <h2>Easy to customize</h2>
              <p>
                중개인과 관리자가 자신의 역할에 맞는 맞춤형 기능을 설정할 수 있도록 유연한 대시보드를 제공합니다.
                사용자는 쉽게 위젯을 추가하거나 삭제할 수 있어 관리가 효율적입니다.
              </p>
            </div>
          </div>
        </section>
      </section>
    </div>
  );
};

export default MainPage;
