import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ApiService } from '../api.service';

@Component({
  selector: 'app-meetups',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <h2>Meetups</h2>
    <div style="margin:8px 0">
      <input [(ngModel)]="q" placeholder="Search meetups" />
      <button (click)="search()">Search</button>
    </div>
    <div *ngIf="loading">Loading…</div>
    <div *ngIf="!loading && items.length === 0">No meetups found.</div>
    <ul>
      <li *ngFor="let m of items" style="margin:8px 0">
        <strong>{{ m.title }}</strong>
        <span *ngIf="m.location"> — {{ m.location }}</span>
        <div style="opacity:.7">{{ m.description }}</div>
        <div>{{ m.isPaid ? ('Paid ' + (m.priceAmount || '')) : 'Free' }}</div>
        <div>{{ m.startTime | date:'medium' }}</div>
      </li>
    </ul>
  `
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
