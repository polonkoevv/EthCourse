export interface Audio {
    id: number;
    title: string;
    cid: string;
    gradient?: string;
    link: string;
    owner_addr: string;
    signature?: string;
    uploaded_at?: string;
  }
  
  export interface Track {
    name: string;
    artist: string;
    link: string; 
    owner_addr: string;
  }