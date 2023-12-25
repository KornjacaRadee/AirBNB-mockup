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
  startDate: Date;
  endDate: Date;
  amenities: string[];

  constructor(
    id: string,
    owner: { id: string },
    name: string,
    location: string,
    minGuestNum: number,
    maxGuestNum: number,
    startDate: Date,
    endDate: Date,
    amenities: string[]
  ) {
    this.id = id;
    this.owner = owner;
    this.name = name;
    this.location = location;
    this.minGuestNum = minGuestNum;
    this.maxGuestNum = maxGuestNum;
    this.startDate = startDate;
    this.endDate = endDate;
    this.amenities = amenities;
  }
}
