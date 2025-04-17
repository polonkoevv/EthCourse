// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract AudioChain {
    struct Audio {
        uint256 id;
        string title;
        string artist;
        string ipfsHash;
        uint256 price;
        address payable owner;
        bool isForSale;
    }
    
    uint256 private nextAudioId = 1;
    mapping(uint256 => Audio) public audios;
    mapping(address => uint256[]) public userAudios;
    mapping(address => mapping(uint256 => bool)) public audioAccess;
    
    uint256 public platformFee = 10; // 10% комиссия платформы
    
    event AudioPublished(uint256 indexed id, address indexed owner, string ipfsHash, uint256 price);
    event AudioPurchased(uint256 indexed id, address indexed buyer, address indexed seller, uint256 amount);
    
    // Публикация аудио
    function publishAudio(string memory _title, string memory _artist, string memory _ipfsHash, uint256 _price) public returns (uint256) {
        uint256 audioId = nextAudioId;
        
        audios[audioId] = Audio({
            id: audioId,
            title: _title,
            artist: _artist,
            ipfsHash: _ipfsHash,
            price: _price,
            owner: payable(msg.sender),
            isForSale: true
        });
        
        userAudios[msg.sender].push(audioId);
        audioAccess[msg.sender][audioId] = true; // Владелец имеет доступ
        
        nextAudioId++;
        
        emit AudioPublished(audioId, msg.sender, _ipfsHash, _price);
        return audioId;
    }
    
    // Покупка аудио
    function purchaseAudio(uint256 _audioId) public payable {
        Audio storage audio = audios[_audioId];
        require(audio.id > 0, "Audio does not exist");
        require(audio.isForSale, "Audio is not for sale");
        require(msg.value >= audio.price, "Insufficient payment");
        require(msg.sender != audio.owner, "Cannot buy your own audio");
        
        uint256 platformAmount = (msg.value * platformFee) / 100;
        uint256 ownerAmount = msg.value - platformAmount;
        
        audio.owner.transfer(ownerAmount);
        // Платформенная комиссия остаётся в контракте
        
        audioAccess[msg.sender][_audioId] = true;
        
        emit AudioPurchased(_audioId, msg.sender, audio.owner, msg.value);
    }
    
    // Проверка доступа к аудио
    function hasAccess(address _user, uint256 _audioId) public view returns (bool) {
        return audioAccess[_user][_audioId];
    }
    
    // Получение всех опубликованных аудио
    function getAudioCount() public view returns (uint256) {
        return nextAudioId - 1;
    }
    
    // Вывод платформенной комиссии (для администратора)
    function withdrawPlatformFees(address payable _recipient) public {
        // Здесь должна быть проверка роли администратора
        _recipient.transfer(address(this).balance);
    }
}