// register.component.ts
import { Component } from '@angular/core';
import { AuthService } from '../services/auth.service';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
})
export class RegisterComponent {
  constructor(private authService: AuthService) {}

  register(user: any): void {
    this.authService.register(user).subscribe(
      (response) => {
        // Obrada uspešne registracije
      },
      (error) => {
        // Obrada greške pri registraciji
      }
    );
  }
}
