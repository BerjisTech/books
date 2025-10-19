import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ApiService } from '../api.service';

@Component({
  selector: 'app-orders',
  standalone: true,
  imports: [CommonModule],
  template: `
    <h2>My Orders</h2>
    <div style="font-size:12px; opacity:.8; margin:8px 0;">Make sure a dev user ID is set on the Marketplace page in development.</div>
    <div *ngIf="loading">Loading…</div>
    <div *ngIf="!loading && items.length === 0">No purchases yet.</div>
    <table *ngIf="items.length > 0" style="border-collapse: collapse; width: 100%;">
      <tr style="text-align:left; border-bottom: 1px solid #e5e7eb22;">
        <th style="padding:8px">Title</th>
        <th style="padding:8px">Kind</th>
        <th style="padding:8px">Status</th>
        <th style="padding:8px">Price</th>
        <th style="padding:8px">Created</th>
      </tr>
      <tr *ngFor="let p of items" style="border-bottom: 1px solid #e5e7eb11;">
        <td style="padding:8px">{{ p.title || p.bookId }}</td>
        <td style="padding:8px">{{ p.kind }}</td>
        <td style="padding:8px">{{ p.status }}</td>
        <td style="padding:8px">{{ p.pricePaid || '-' }} {{ p.currency || '' }}</td>
        <td style="padding:8px">{{ p.createdAt | date:'short' }}</td>
      </tr>
    </table>
  `
})
export class OrdersPageComponent {
  loading = false;
  items: any[] = [];
  constructor(private api: ApiService) { this.refresh(); }
  refresh() {
    this.loading = true;
    this.api.listMyPurchases().subscribe({
      next: res => { this.items = res.data || []; this.loading = false; },
      error: () => { this.items = []; this.loading = false; }
    });
  }
}

