// auth.service.ts
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class AuthService { // Postavite odgovarajuću adresu backend-a

  constructor(private http: HttpClient) {}

}
