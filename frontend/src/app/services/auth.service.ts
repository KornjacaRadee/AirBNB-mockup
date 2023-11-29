import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { ConfigService } from './config.service';
import { tap } from 'rxjs';
import * as jwt_decode_ from 'jwt-decode';
import { Router } from '@angular/router';
import { catchError} from 'rxjs/operators';
import { User } from '../model/user';




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
  whoami(): Observable<User> {
    const userId = this.getCurrentUserId();

    if (!userId) {
      // Throw an error observable with a custom message
      return throwError('User ID not found in the token.');
    }

    return this.http.get<User>(`${this.configService._whoami_url}${userId}`).pipe(
      catchError((error) => {
        // Handle HTTP errors or return a default user object
        console.error('Error fetching user information:', error);
        return throwError('Failed to fetch user information.');
      })
    );
  }

  private getCurrentUserId(): string | null {
    // Logic to get the current user's ID from the token or elsewhere
    const token = this.getAuthToken();
    if (token) {
      const decodedToken: any = jwt_decode_ as any;
      console.log(decodedToken.user_id) // Type assertion to any
      return decodedToken(token).user_id;
       // Adjust this based on your JWT payload
    }
    return null;
  }

}
