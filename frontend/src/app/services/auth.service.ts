import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ConfigService } from './config.service';

interface User {
  username: string;
  email: string;
  password?: string;
  firstName?: string;
  address?: string;

  // Dodajte druge atribute korisnika prema potrebi
}

interface LoginCredentials {
  email: string;
  password: string;
}

@Injectable({
  providedIn: 'root',
})
export class AuthService {

  constructor(private http: HttpClient,
    private configService:ConfigService) {


  }

  register(user: User): Observable<any> {
    return this.http.post(`${this.configService._register_url}`, user);
  }

  login(credentials: LoginCredentials): Observable<any> {
    return this.http.post(`${this.configService._login_url}`, credentials);
  }
}
