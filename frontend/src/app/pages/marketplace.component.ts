import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ApiService } from '../api.service';
import { PaymentPanelComponent } from '../components/payment-panel/payment-panel.component';
import { interval, Subscription } from 'rxjs';

@Component({
  selector: 'app-marketplace',
  standalone: true,
  imports: [CommonModule, FormsModule, PaymentPanelComponent],
  template: `
    <h2>Marketplace</h2>
    <div style="margin:8px 0">
      <input [(ngModel)]="q" placeholder="Search books" />
      <button (click)="search()">Search</button>
    </div>
    <div style="margin:8px 0; font-size:12px; opacity:.8">
      Dev user ID: <input [(ngModel)]="devUserId" placeholder="dev user id" style="width:360px" />
      <button (click)="saveDevUser()">Save</button>
    </div>
    <div *ngIf="loading">Loading…</div>
    <div *ngIf="!loading && items.length === 0">No books found.</div>
    <ul>
      <li *ngFor="let b of items" style="margin:8px 0">
        <strong>{{ b.title }}</strong>
        <span *ngIf="b.subtitle"> - {{ b.subtitle }}</span>
        <div style="opacity:.7">{{ b.description }}</div>
        <div *ngIf="b.priceAmount">Price: {{ b.priceAmount }} {{ b.priceCurrency }}</div>
        <div>
          <button (click)="buy(b)">Buy (M-Pesa)</button>
        </div>
        <div *ngIf="selected && selected.id === b.id && paymentDetails" style="margin-top:10px;">
          <app-payment-panel [details]="paymentDetails"></app-payment-panel>
        </div>
      </li>
    </ul>
  `
})
export class MarketplacePageComponent {
  items: any[] = [];
  loading = false;
  q = '';
  selected: any = null;
  paymentDetails: any = null;
  pollSub?: Subscription;
  devUserId = localStorage.getItem('devUserId') || '';
  constructor(private api: ApiService) { this.search(); }
  search() {
    this.loading = true;
    this.api.getPublicBooks(this.q).subscribe({
      next: res => { this.items = res.data || []; this.loading = false; },
      error: () => { this.items = []; this.loading = false; }
    });
  }
  buy(b: any) {
    this.selected = b;
    this.paymentDetails = null;
    this.api.createPurchase(b.id, 'ebook', 'mpesa').subscribe(res => {
      if (res.status === 402) {
        this.paymentDetails = res.body.data;
        const intent = this.paymentDetails?.payment?.intent;
        if (intent?.id) {
          // start polling Core API for intent status
          this.pollSub?.unsubscribe();
          this.pollSub = interval(5000).subscribe(() => this.refreshIntent());
        }
      } else if (res.status >= 200 && res.status < 300) {
        alert('Purchase completed');
      } else {
        alert('Failed to start purchase');
      }
    });
  }
  refreshIntent() {
    const intent = this.paymentDetails?.payment?.intent;
    if (!intent?.id) return;
    this.api.getPaymentIntent(intent.id).subscribe(resp => {
      const st = (resp.data?.status || '').toLowerCase();
      if (st && st !== 'pending' && st !== 'processing') {
        // stop polling
        this.pollSub?.unsubscribe();
        // Confirm purchase in Books so Orders reflect completion
        const purchaseId = this.paymentDetails?.purchaseId;
        if (purchaseId) {
          this.api.confirmPurchase(purchaseId).subscribe(() => {
            alert('Payment ' + st + ' — purchase confirmed');
          });
        } else {
          alert('Payment ' + st);
        }
      }
    });
  }
  saveDevUser() {
    if (this.devUserId) localStorage.setItem('devUserId', this.devUserId); else localStorage.removeItem('devUserId');
    alert('Saved');
  }
}
