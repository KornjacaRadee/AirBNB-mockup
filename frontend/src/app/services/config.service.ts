import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class ConfigService {


  _api_url: string;
  _login_url: string;
  _register_url:string;
  _accomodations_url:string;
  _reservations_url:string;


  constructor() {
    this._api_url = 'https://localhost/';
    this._login_url =this._api_url + 'auth/login';
    this._register_url = this._api_url + 'auth/register';
    this._accomodations_url = this._api_url + 'accommodations';
    this._reservations_url = this._api_url + 'reservations';
    


  }
}
