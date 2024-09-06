export interface Raid {
  ID: number;
  Name: string;
  Slug: string;
  ShortName: string;
  Expansion: string;
  MediaURL: string;
  Icon: string;
  Starts: {
    US: string;
    EU: string;
    TW: string;
    KR: string;
    CN: string;
  };
  Ends: {
    US: string;
    EU: string;
    TW: string;
    KR: string;
    CN: string;
  };
  Encounters: {
    ID: number;
    Slug: string;
    Name: string;
  }[];
}
