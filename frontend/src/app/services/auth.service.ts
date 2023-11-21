import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

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
  private apiUrl = 'http://localhost:8080';

  constructor(private http: HttpClient) {}

  register(user: User): Observable<any> {
    return this.http.post(`${this.apiUrl}/register`, user);
  }

  login(credentials: LoginCredentials): Observable<any> {
    return this.http.post(`${this.apiUrl}/login`, credentials);
  }
}
