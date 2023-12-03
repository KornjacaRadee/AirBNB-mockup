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
  _recovery_url:string;
  _validatetoken_url:string;
  _updatenewpassword_url: string;
  _getuserbyid_url: string;
  _createaccom_url: string;

  constructor() {
    this._api_url = 'https://localhost/';
    this._login_url =this._api_url + 'auth/login';
    this._register_url = this._api_url + 'auth/register';
    this._accomodations_url = this._api_url + 'accommodations';
    this._reservations_url = this._api_url + 'reservations';
    this._recovery_url = this._api_url + 'auth/password-recovery';
    this._validatetoken_url = this._api_url + 'auth/reset';
    this._updatenewpassword_url = this._api_url + 'auth/update';
    this._getuserbyid_url = this._api_url + 'auth/users';


    this._createaccom_url = this._api_url + 'accommodations/new';




  }
}
