import { Component, Input } from '@angular/core';
import { Book } from '../../interfaces/book';
import { BookCardComponent } from "../book-card/book-card.component";
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-books-list',
  standalone: true,
  templateUrl: './books-list.component.html',
  styleUrl: './books-list.component.scss',
  imports: [CommonModule, BookCardComponent]
})
export class BooksListComponent {
  // Array of books
  @Input() books: Book[] = [];

  @Input() pageTitle: string = 'Book Recomendations';
}
