export class User {
  id: string;
  firstName: string;
  lastName: string | null;
  email: string;
  password: string;
  address: string | null;
  createdOn: string;
  updatedOn: string;
  deletedOn: string;
  roles: string[];

  constructor(
    id: string,
    firstName: string,
    lastName: string | null,
    email: string,
    password: string,
    address: string | null,
    createdOn: string,
    updatedOn: string,
    deletedOn: string,
    roles: string[]
  ) {
    this.id = id;
    this.firstName = firstName;
    this.lastName = lastName;
    this.email = email;
    this.password = password;
    this.address = address;
    this.createdOn = createdOn;
    this.updatedOn = updatedOn;
    this.deletedOn = deletedOn;
    this.roles = roles;
  }
}
