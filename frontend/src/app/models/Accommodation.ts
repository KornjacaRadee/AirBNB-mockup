// accommodation.model.ts

export class Accommodation {
  id: string;
  owner: {
    id: string;
  };
  name: string;
  location: string;
  minGuestNum: number;
  maxGuestNum: number;
  amenities: string[];

  constructor(
    id: string,
    owner: { id: string },
    name: string,
    location: string,
    minGuestNum: number,
    maxGuestNum: number,
    amenities: string[]
  ) {
    this.id = id;
    this.owner = owner;
    this.name = name;
    this.location = location;
    this.minGuestNum = minGuestNum;
    this.maxGuestNum = maxGuestNum;
    this.amenities = amenities;
  }
}
