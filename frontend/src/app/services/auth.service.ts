import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ConfigService } from './config.service';
import { tap } from 'rxjs';
import * as jwt_decode_ from 'jwt-decode';
import { Router } from '@angular/router';


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

  private tokenKey = 'authToken';

  constructor(
    private http: HttpClient,
    private configService:ConfigService,
    private router: Router
    ) {


  }

  register(user: User): Observable<any> {
    return this.http.post(`${this.configService._register_url}`, user);
  }

  // login(credentials: LoginCredentials): Observable<any> {
  //   return this.http.post(`${this.configService._login_url}`, credentials);
  // }

  login(credentials: LoginCredentials): Observable<any> {
    return this.http.post(`${this.configService._login_url}`, credentials).pipe(
      tap((response: any) => {
        const token = response.token;
        if (token) {
          // Store the token in localStorage
          localStorage.setItem(this.tokenKey, token);
        }
      })
    );
  }
  logout() {
    console.log('Logout method called');
    localStorage.removeItem(this.tokenKey);
    console.log('Navigating to lgin page');
    this.router.navigate(['/login']);
  }

  getAuthToken(): string | null {
    return localStorage.getItem(this.tokenKey);
  }

  getUserRole(): string | null {
    const token = this.getAuthToken();
    if (token) {
      const decodedToken: any = jwt_decode_ as any; // Type assertion to any
      return decodedToken(token).roles; // Adjust this based on your JWT payload
    }
    return null;
  }

  isAuthenticated(): boolean {
    // Check if the user is authenticated based on the presence of the token
    return !!this.getAuthToken();
  }

  recovery(email: any): Observable<any>{
    return this.http.post(`${this.configService._recovery_url}`, email);
  }
  validateToken(token: string): Observable<any>{
    return this.http.get(`${this.configService._validatetoken_url}?token=${token}`);
  }
  setNewPassword(token: string,password: any): Observable<any>{
    return this.http.post(`${this.configService._updatenewpassword_url}?token=${token}`, password);
  }

}
