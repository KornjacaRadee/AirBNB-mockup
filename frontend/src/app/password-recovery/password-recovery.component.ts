import { Component } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-password-recovery',
  templateUrl: './password-recovery.component.html',
  styleUrls: ['./password-recovery.component.css']
})
export class PasswordRecoveryComponent {
  email: string = '';

  constructor(private authService: AuthService, private router: Router) {}

  sendRecoveryEmail() {
    const payload = { email: this.email };

    this.authService.recovery(payload).subscribe(
      (response) => {
        alert("If the user with the email exists, a recovery token has been sent. Check your inbox.");
        this.router.navigate(['/new-password']);
        console.log('Recovery email sent successfully', response);
      },
      (error) => {
        console.error('Failed to send recovery email', error);
      }
    );
  }
}
