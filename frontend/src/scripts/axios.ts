import axios, { type AxiosInstance } from 'axios';

const baseURL = 'http://localhost:8000';

class AxiosEntity {
  private api: AxiosInstance;
  constructor() {
    this.api = axios.create({
      baseURL: baseURL,
    });
  }

  async GetAllMusic() {
    return this.api.get(baseURL + '/music');
  }

  async GetMusicById(id: string) {
    return this.api.get(baseURL + '/music/' + id);
  }

  async UploadMusic(music: Blob, trackTitle: string) {
    
    const formData = new FormData();
    formData.append('file', music);
    formData.append('title', trackTitle);
    return this.api.post(baseURL + '/upload', formData);
  }
  
  
}

export default new AxiosEntity();