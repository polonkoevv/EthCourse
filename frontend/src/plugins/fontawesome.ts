import { library } from '@fortawesome/fontawesome-svg-core';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';

// Иконки из solid набора
import { 
  faMusic, 
  faPlay, 
  faPause, 
  faCloudArrowUp, 
  faShareNodes, 
  faEllipsis, 
  faBackwardStep, 
  faForwardStep, 
  faArrowRightArrowLeft,
  faUpload,
  faXmark,
  // Новые иконки для работы с кошельком
  faUser,
  faWallet,
  faArrowRight,
  faCheck,
  faArrowLeft
} from '@fortawesome/free-solid-svg-icons';

// Иконки из regular набора
import { faMoon, faSun } from '@fortawesome/free-regular-svg-icons';

// Иконки из brands набора
import { faTwitter, faDiscord, faGithub } from '@fortawesome/free-brands-svg-icons';

// Добавляем иконки в библиотеку
library.add(
  // Solid
  faMusic, 
  faPlay, 
  faPause, 
  faCloudArrowUp, 
  faShareNodes, 
  faEllipsis, 
  faBackwardStep, 
  faForwardStep, 
  faArrowRightArrowLeft,
  faUpload,
  faXmark,
  // Новые иконки
  faUser,
  faWallet,
  faArrowRight,
  faCheck,
  faArrowLeft,
  
  // Regular
  faMoon, 
  faSun,
  
  // Brands
  faTwitter, 
  faDiscord, 
  faGithub
);

export { FontAwesomeIcon };