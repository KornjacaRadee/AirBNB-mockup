import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { ActivatedRoute } from '@angular/router';
import {  HttpHeaders } from '@angular/common/http';
import { AccomodationService } from '../services/accomodation/accomodation.service';
import { ToastrService } from 'ngx-toastr';
import { AuthService } from '../services/auth.service';

interface PictureData {
  id: number;
  accommodationID: string;
  data: File;
  base64Data?: string;
}

@Component({
  selector: 'app-add-picutres-accomm',
  templateUrl: './add-picutres-accomm.component.html',
  styleUrls: ['./add-picutres-accomm.component.css']
})


export class AddPicutresAccommComponent implements OnInit {
  accommID: string = "";
  picturesData: PictureData[] = [];

  ngOnInit(): void {
    this.route.queryParams.subscribe(params => {
      this.accommID = params['id'];

    });
  }

  constructor( private toastr: ToastrService,private http: HttpClient,private router: Router,private route: ActivatedRoute, private accommodationService: AccomodationService,private authService: AuthService) {}


  onFileSelected(event: any): void {
    const files: FileList = event.target.files;

    for (let i = 0; i < files.length && this.picturesData.length < 10; i++) {
      if (this.picturesData.length < 10) {
        const picture: PictureData = {
          id: i,
          accommodationID: this.accommID,
          data: files[i]
        };

        this.encodeFileToBase64(files[i])
          .then(base64Data => {
            const parts = base64Data.split(',');

            base64Data = parts.length === 2 ? parts[1] : '';
            picture.base64Data = base64Data;
          })
          .catch(error => {
            console.error('Failed to encode file:', error);
          });

        this.picturesData.push(picture);
      } else {
        break;
      }
    }
  }

  onUpload(): void {
    const token = this.authService.getAuthToken();
    const headers = new HttpHeaders().set('Authorization', `Bearer ${token}`);

    if (this.picturesData.length > 0) {
      const jsonData = this.picturesData.map(picture => ({
        id: picture.id.toString(), // Convert id to string
        accommodation_id: this.accommID, // Change to underscore notation
        data: picture.base64Data // Use Base64 data
      }));

      this.accommodationService.addAccommodationPictures(jsonData)
        .subscribe(
          (response) => {
            this.toastr.success('Pictures added successfully');
            console.log('Accommodation created successfully', response);
          },
          (error) => {
            this.toastr.error('Failed to add pictures!');
            console.error('Failed to create accommodation', error);
          }
        );
    }
  }

  encodeFileToBase64(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = () => {
        resolve(reader.result?.toString() || '');
      };
      reader.onerror = reject;
      reader.readAsDataURL(file);
    });
  }

}
