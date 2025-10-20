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
  templateUrl: './marketplace.component.html'
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
