#include <iostream>
#include <vector>
#include <ctime>
#include <sstream>
#include <iomanip>
#include <openssl/sha.h>  // OpenSSL 라이브러리를 사용하여 SHA256 해시를 계산합니다.

// 블록 구조체 정의
struct Block {
    int index;
    std::string timestamp;
    std::string data;   // 이 부분에 매물 정보와 같은 데이터를 추가할 수 있습니다.
    std::string prevHash;
    std::string hash;

    Block(int idx, std::string ts, std::string d, std::string prevH) {
        index = idx;
        timestamp = ts;
        data = d;
        prevHash = prevH;
        hash = calculateHash();
    }

    // 블록의 해시를 계산하는 함수
    std::string calculateHash() const {
        std::stringstream ss;
        ss << index << timestamp << data << prevHash;
        std::string input = ss.str();
        unsigned char hash[SHA256_DIGEST_LENGTH];
        SHA256_CTX sha256;
        SHA256_Init(&sha256);
        SHA256_Update(&sha256, input.c_str(), input.length());
        SHA256_Final(hash, &sha256);

        std::stringstream hashStringStream;
        for (int i = 0; i < SHA256_DIGEST_LENGTH; i++) {
            hashStringStream << std::setw(2) << std::setfill('0') << std::hex << (int)hash[i];
        }
        return hashStringStream.str();
    }
};

// 블록체인 클래스 정의
class Blockchain {
private:
    std::vector<Block> chain;

public:
    // 생성자: 제네시스 블록을 생성
    Blockchain() {
        // 제네시스 블록을 수동으로 추가
        chain.push_back(createGenesisBlock());
    }

    // 제네시스 블록 생성 함수
    Block createGenesisBlock() {
        return Block(0, "2025-03-21", "Genesis Block", "0");
    }

    // 새 블록 추가 함수
    void addBlock(std::string data) {
        int index = chain.size();
        std::string timestamp = getCurrentTime();
        std::string prevHash = chain.back().hash;
        Block newBlock(index, timestamp, data, prevHash);
        chain.push_back(newBlock);
    }

    // 블록체인 출력 함수
    void printBlockchain() {
        for (const Block& block : chain) {
            std::cout << "Block #" << block.index << std::endl;
            std::cout << "Timestamp: " << block.timestamp << std::endl;
            std::cout << "Data: " << block.data << std::endl;
            std::cout << "Previous Hash: " << block.prevHash << std::endl;
            std::cout << "Hash: " << block.hash << std::endl;
            std::cout << std::endl;
        }
    }

    // 현재 시간 반환
    std::string getCurrentTime() {
        time_t now = time(0);
        tm* ltm = localtime(&now);
        std::stringstream ss;
        ss << 1900 + ltm->tm_year << '-' 
           << std::setw(2) << std::setfill('0') << 1 + ltm->tm_mon << '-'
           << std::setw(2) << std::setfill('0') << ltm->tm_mday;
        return ss.str();
    }
};

int main() {
    Blockchain myBlockchain;

    // 블록체인에 매물 정보를 추가 (예시: 부동산 매물)
    myBlockchain.addBlock("매물 1: 서울 강남구 아파트, 가격: 10억");
    myBlockchain.addBlock("매물 2: 서울 마포구 오피스텔, 가격: 5억");

    // 블록체인 출력
    myBlockchain.printBlockchain();

    return 0;
}
