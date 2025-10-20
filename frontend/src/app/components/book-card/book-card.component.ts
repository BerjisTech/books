import { Component, Input } from '@angular/core';
import { Book } from '../../interfaces/book';
import { CommonModule } from '@angular/common';
import { RouterLink } from "@angular/router";

@Component({
  selector: 'app-book-card',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './book-card.component.html',
  styleUrl: './book-card.component.scss'
})
export class BookCardComponent {
  @Input() book: Book|null = null;
}
