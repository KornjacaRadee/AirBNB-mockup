import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class ConfigService {


  _api_url: string;
  _login_url: string;
  _register_url:string;

  constructor() {
    this._api_url = 'http://localhost:8080/';
    this._login_url =this._api_url + '/login';
    this._register_url = this._api_url + '/register';
    

  }
}
