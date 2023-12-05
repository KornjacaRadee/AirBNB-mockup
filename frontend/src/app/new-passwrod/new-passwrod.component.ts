import { Component } from '@angular/core';
import { HttpClientModule } from '@angular/common/http'
import { AuthService } from '../services/auth.service';
import { Router } from '@angular/router';
@Component({
  selector: 'app-new-passwrod',
  templateUrl: './new-passwrod.component.html',
  styleUrls: ['./new-passwrod.component.css']
})
export class NewPasswrodComponent {
  recoveryToken: string = '';
  newPassword: string = '';
  repeatPassword: string = '';
  showNewPasswordForm: boolean = false;

  constructor(private authService: AuthService, private router: Router) {}

  validateToken() {
    const payload = { token: this.recoveryToken };

    this.authService.validateToken(this.recoveryToken).subscribe(
      (response) => {
        console.log(response)
        if (response && response.status === 'success') {
          this.showNewPasswordForm = true;
        } else {
          alert('Invalid recovery token. Please try again.');
          // Optionally, you can reset the form or take additional actions
        }
      },
      (error) => {
        console.error('Error validating token', error);
        // Handle other errors
      }
    );
  }


  setNewPassword() {
    if (this.newPassword === this.repeatPassword) {
      const payload = { password: this.newPassword };

      this.authService.setNewPassword(this.recoveryToken, payload).subscribe(
        (response) => {
          alert('Password set successfully!');
          this.router.navigate(['/login']); // Redirect to login page or any other page
        },
        (error) => {
          console.error('Error setting new password', error);
          // Handle error
        }
      );
    } else {
      alert('Passwords do not match. Please make sure the passwords match.');
    }
  }
}
