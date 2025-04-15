export interface Transaction {
    type: string;
    title: string;
    address: string;
    amount: string;
    icon: string | string[]; // Изменено для поддержки формата FontAwesome
    iconColor: string;
  }