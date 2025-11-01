import { Component, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ApiService } from '../api.service';
import { PaymentPanelComponent } from '../components/payment-panel/payment-panel.component';
import { interval, Subscription } from 'rxjs';
import { Router } from '@angular/router';

type Tone = 'info' | 'success' | 'error';

@Component({
  selector: 'app-marketplace',
  standalone: true,
  imports: [CommonModule, FormsModule, PaymentPanelComponent],
  templateUrl: './marketplace.component.html'
})
export class MarketplacePageComponent implements OnDestroy {
  items: any[] = [];
  loading = false;
  q = '';
  selected: any = null;
  paymentDetails: any = null;
  pollSub?: Subscription;
  devUserId = localStorage.getItem('devUserId') || '';
  notice: { text: string; tone: Tone } | null = null;
  constructor(private api: ApiService, private router: Router) { this.search(); }

  ngOnDestroy(): void {
    this.pollSub?.unsubscribe();
  }
  private setNotice(text: string, tone: Tone = 'info') {
    this.notice = { text, tone };
  }
  clearNotice() {
    this.notice = null;
  }

  search() {
    this.loading = true;
    this.clearNotice();
    this.api.getPublicBooks(this.q).subscribe({
      next: res => {
        const data = res.data || [];
        this.items = data.map((b: any) => ({
          ...b,
          displayCover: b.coverImageUrl,
          displayAuthor: b.authorName,
          slugOrId: b.slug || b.id
        }));
        this.loading = false;
      },
      error: () => { this.items = []; this.loading = false; this.setNotice('We could not load the marketplace right now.', 'error'); }
    });
  }
  coverFor(book: any) {
    if (book?.displayCover) return book.displayCover;
    const fallback = this.items.length ? this.items[0].displayCover : null;
    return fallback || 'https://www.rockingbookcovers.com/wp-content/uploads/2023/03/Dont-Go-There.jpg';
  }
  view(book: any) {
    const target = book?.slugOrId || book?.id;
    if (target) {
      this.router.navigate(['/book', target]);
    }
  }
  buy(b: any) {
    this.selected = b;
    this.paymentDetails = null;
    this.setNotice('Preparing checkout…', 'info');
    this.api.createPurchase(b.id, 'ebook', 'mpesa').subscribe(res => {
      if (res.status === 402) {
        this.paymentDetails = res.body.data;
        const intent = this.paymentDetails?.payment?.intent;
        if (intent?.id) {
          // start polling Core API for intent status
          this.pollSub?.unsubscribe();
          this.pollSub = interval(5000).subscribe(() => this.refreshIntent());
        }
        this.setNotice('Payment initiated. Complete it using the instructions below.', 'info');
      } else if (res.status >= 200 && res.status < 300) {
        this.setNotice('Purchase completed. You can find it in your library.', 'success');
      } else {
        this.setNotice('Unable to start that purchase. Please try again.', 'error');
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
            this.setNotice(`Payment ${st}. Purchase confirmed.`, st === 'succeeded' ? 'success' : 'info');
          });
        } else {
          this.setNotice(`Payment ${st}.`, st === 'succeeded' ? 'success' : 'info');
        }
      }
    });
  }
  saveDevUser() {
    if (this.devUserId) localStorage.setItem('devUserId', this.devUserId); else localStorage.removeItem('devUserId');
    this.setNotice('Developer user preference saved for this browser.', 'success');
  }
}
