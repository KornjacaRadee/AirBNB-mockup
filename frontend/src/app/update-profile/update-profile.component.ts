// update-profile.component.ts
import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { ProfilesService } from '../services/profile/profiles.service';

@Component({
  selector: 'app-update-profile',
  templateUrl: './update-profile.component.html',
  styleUrls: ['./update-profile.component.css'],
})
export class UpdateProfileComponent implements OnInit {
  profile: any = {};
  email: string = '';
  updatedProfile: any = {
    address: '',
    email: '',
    username: '',
  };

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private profileService: ProfilesService
  ) {}

  ngOnInit(): void {
    this.route.queryParams.subscribe((params) => {
      this.email = params['email'];
      this.loadProfileDetails();
    });
  }

  loadProfileDetails() {
    this.profileService.getProfileByEmail(this.email).subscribe(
      (response) => {
        // Map the response to the 'user' property
        this.profile = response;
        console.log(this.profile);
      },
      (error) => {
        console.error('Error fetching user details', error);
      }
    );
  }

  saveProfile(): void {
    this.profileService
      .updateProfileByEmail(this.email, this.profile)
      .subscribe(
        () => {
          this.router.navigate(['/profiles']);
        },
        (error: any) => {
          console.error('Error:', error);
        }
      );
  }
}
