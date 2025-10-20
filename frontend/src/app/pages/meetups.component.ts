import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ApiService } from '../api.service';

@Component({
  selector: 'app-meetups',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './meetups.component.html'
})
export class MeetupsPageComponent {
  items: any[] = [];
  loading = false;
  q = '';
  constructor(private api: ApiService) { this.search(); }
  search() {
    this.loading = true;
    this.api.getPublicMeetups(this.q).subscribe({
      next: res => { this.items = res.data || []; this.loading = false; },
      error: () => { this.items = []; this.loading = false; }
    });
  }
}
